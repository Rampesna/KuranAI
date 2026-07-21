import SwiftUI

struct SourceSelector: View {
    let language: Language
    @Binding var selected: SourceFilter

    var body: some View {
        HStack(spacing: 12) {
            ForEach(SourceFilter.allCases, id: \.self) { filter in
                Button {
                    selected = filter
                } label: {
                    Text(filter.title(for: language))
                        .font(.system(size: 13, weight: .semibold))
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 12)
                        .background(
                            RoundedRectangle(cornerRadius: 18)
                                .fill(selected == filter ? Theme.neon : Theme.pillGradient)
                        )
                        .foregroundColor(selected == filter ? .black : Theme.textPrimary)
                }
            }
        }
    }
}
